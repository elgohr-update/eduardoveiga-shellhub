<template>
  <fragment>
    <v-tooltip bottom>
      <template #activator="{ on }">
        <v-icon
          v-on="on"
          @click="open()"
        >
          mdi-console
        </v-icon>
      </template>
      <span data-test="console-tooltip">
        Terminal
      </span>
    </v-tooltip>

    <v-dialog
      v-model="show"
      max-width="1024px"
    >
      <v-card data-test="terminal-dialog">
        <v-toolbar
          dark
          color="primary"
        >
          <v-btn
            icon
            dark
            @click="close()"
          >
            <v-icon>close</v-icon>
          </v-btn>
          <v-toolbar-title>Terminal</v-toolbar-title>

          <v-spacer />
        </v-toolbar>

        <v-card
          v-if="showLoginForm"
          class="ma-0 px-6 py-4"
          outlined
        >
          <v-tabs
            centered
          >
            <v-tab
              v-for="tab in tabs"
              :key="tab"
              :data-test="tab+'-tab'"
              @click="resetFieldValidation"
            >
              {{ tab }}
            </v-tab>

            <v-tab-item>
              <v-card
                flat
              >
                <v-form
                  ref="form"
                  v-model="valid"
                  lazy-validation
                  @submit.prevent="connectWithPassword()"
                >
                  <v-text-field
                    ref="username"
                    v-model="username"
                    label="Username"
                    autofocus
                    :rules="[rules.required]"
                    :validate-on-blur="true"
                    data-test="username-field"
                  />

                  <v-text-field
                    ref="passwd"
                    v-model="passwd"
                    label="Password"
                    type="password"
                    :rules="[rules.required]"
                    :validate-on-blur="true"
                    data-test="passwd-field"
                  />

                  <v-card-actions>
                    <v-spacer />
                    <v-btn
                      type="submit"
                      color="primary"
                      class="mt-4"
                      data-test="connect-btn"
                    >
                      Connect
                    </v-btn>
                  </v-card-actions>
                </v-form>
              </v-card>
            </v-tab-item>

            <v-tab-item>
              <v-card
                flat
              >
                <v-form
                  ref="form"
                  v-model="valid"
                  lazy-validation
                  @submit.prevent="connectWithPrivateKey()"
                >
                  <v-text-field
                    ref="username"
                    v-model="username"
                    label="Username"
                    autofocus
                    :rules="[rules.required]"
                    :validate-on-blur="true"
                    data-test="username2-field"
                  />

                  <v-select
                    v-model="privateKey"
                    :items="getListPrivateKeys"
                    item-text="name"
                    item-value="data"
                    label="Private Keys"
                    data-test="privatekeys-select"
                  />

                  <v-card-actions>
                    <v-spacer />
                    <v-btn
                      type="submit"
                      color="primary"
                      class="mt-4"
                      data-test="connect2-btn"
                    >
                      Connect
                    </v-btn>
                  </v-card-actions>
                </v-form>
              </v-card>
            </v-tab-item>
          </v-tabs>
        </v-card>
        <div ref="terminal" />
      </v-card>
    </v-dialog>
  </fragment>
</template>

<script>

import { Terminal } from 'xterm';
import { AttachAddon } from 'xterm-addon-attach';
import { FitAddon } from 'xterm-addon-fit';

import 'xterm/css/xterm.css';

import RSAKey from 'node-rsa';
import { parsePrivateKey } from '@/sshpk';

export default {
  name: 'TerminalDialog',

  props: {
    uid: {
      type: String,
      required: true,
    },
  },

  data() {
    return {
      username: '',
      passwd: '',
      showLoginForm: true,
      valid: true,
      privateKey: '',
      rules: {
        required: (value) => !!value || 'Required',
      },
      tabs: ['Password', 'PublicKey'],
    };
  },

  computed: {
    webTermDimensions() {
      return {
        cols: this.xterm.cols,
        rows: this.xterm.rows,
      };
    },

    show: {
      get() {
        return this.$store.getters['modals/terminal'] === this.$props.uid;
      },

      set(value) {
        if (value) {
          this.$store.dispatch('modals/toggleTerminal', this.$props.uid);
        } else {
          this.$store.dispatch('modals/toggleTerminal', '');
        }
      },
    },

    getListPrivateKeys() {
      return this.$store.getters['privatekeys/list'];
    },
  },

  watch: {
    show(value) {
      if (!value) {
        if (this.ws) this.ws.close();
        if (this.xterm) this.xterm.dispose();

        this.username = '';
        this.passwd = '';
        this.showLoginForm = true;
      } else {
        requestAnimationFrame(() => {
          this.$refs.username.focus();
        });
      }
    },
  },

  methods: {
    open() {
      this.privateKey = '';

      this.xterm = new Terminal({
        cursorBlink: true,
        fontFamily: 'monospace',
      });

      this.fitAddon = new FitAddon();
      this.xterm.loadAddon(this.fitAddon);

      this.$store.dispatch('modals/toggleTerminal', this.$props.uid);

      if (this.xterm.element) {
        this.xterm.reset();
      }
    },

    close() {
      this.$store.dispatch('modals/toggleTerminal', '');
    },

    connectWithPassword() {
      const passwd = encodeURIComponent(this.passwd);
      this.connect({ passwd });
    },

    encodeURLParams(params) {
      return Object.entries(params).map(([k, v]) => `${k}=${v}`).join('&');
    },

    connectWithPrivateKey() {
      const key = new RSAKey(this.privateKey);
      key.setOptions({ signingScheme: 'pkcs1-sha1' });

      const signature = encodeURIComponent(key.sign(this.username, 'base64'));
      const fingerprint = parsePrivateKey(this.privateKey).fingerprint('md5');

      this.connect({ signature, fingerprint });
    },

    connect(params) {
      if (!this.$refs.form.validate(true)) {
        return;
      }

      this.showLoginForm = false;
      this.$nextTick(() => this.fitAddon.fit());

      if (!this.xterm.element) {
        this.xterm.open(this.$refs.terminal);
      }

      this.fitAddon.fit();
      this.xterm.focus();

      let protocolConnectionURL = '';

      if (window.location.protocol === 'http:') {
        protocolConnectionURL = 'ws';
      } else {
        protocolConnectionURL = 'wss';
      }

      const wsInfo = { user: `${this.username}@${this.$props.uid}`, ...params, ...this.webTermDimensions };
      this.ws = new WebSocket(`${protocolConnectionURL}://${window.location.host}/ws/ssh?${this.encodeURLParams(wsInfo)}`);

      this.ws.onopen = () => {
        this.attachAddon = new AttachAddon(this.ws);
        this.xterm.loadAddon(this.attachAddon);
      };

      this.ws.onclose = () => {
        this.attachAddon.dispose();
      };
    },

    resetFieldValidation() {
      this.$refs.username.reset();
      this.$refs.passwd.reset();
    },
  },
};

</script>
