<template>
  <v-form>
    <v-container>
      <v-row
        align="center"
        justify="center"
        class="mt-4 mb-4"
      >
        <v-col
          sm="8"
        >
          <v-card class="pb-0 elevation-0">
            <div class="d-flex pa-0 align-center">
              <v-spacer />
              <v-spacer />
              <PrivateKeyFormDialog
                :create-key="true"
                action="private"
                data-test="privateKeyFormDialogFirst-component"
              />
            </div>

            <v-data-table
              :headers="headers"
              :items="getListPrivateKeys"
              data-test="dataTable-field"
              :server-items-length="getNumberPrivateKeys"
              hide-default-footer
            >
              <template #[`item.name`]="{ item }">
                {{ item.name }}
              </template>

              <template #[`item.data`]="{ item }">
                {{ convertToFingerprint(item.data) }}
              </template>

              <template #[`item.actions`]="{ item }">
                <PrivateKeyFormDialog
                  :key-object="item"
                  :create-key="false"
                  action="private"
                  data-test="privateKeyFormDialogSecond-component"
                />
                <PrivateKeyDelete
                  :fingerprint="item.data"
                  action="private"
                  data-test="privateKeyDelete-component"
                />
              </template>
            </v-data-table>
          </v-card>
        </v-col>
      </v-row>
    </v-container>

    <v-dialog
      v-model="dialog"
      persistent
      width="500"
    >
      <v-card>
        <v-card-title class="headline grey lighten-2">
          Privacy Policy
        </v-card-title>

        <v-card-text
          class="mt-4"
        >
          The private key is never submitted to ShellHub, it gets stored in your browser’s
          local storage, only the public key gets uploaded and stored by ShellHub.
        </v-card-text>

        <v-divider />
        <v-card-actions
          class="px-6"
        >
          <v-checkbox
            v-model="privatekeyPrivacyPolicy"
            label="Never show this again"
          />
          <v-spacer />
          <v-btn
            color="primary"
            text
            data-test="gotIt-btn"
            @click="accept"
          >
            Got it
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>

<script>

import PrivateKeyFormDialog from '@/components/public_key/KeyFormDialog';
import PrivateKeyDelete from '@/components/public_key/KeyDelete';

import { parsePrivateKey } from '@/sshpk';

export default {
  name: 'SettingPrivateKeys',

  components: {
    PrivateKeyFormDialog,
    PrivateKeyDelete,
  },

  data() {
    return {
      pagination: {},
      dialog: true,
      privatekeyPrivacyPolicy: false,
      headers: [
        {
          text: 'Name',
          value: 'name',
          align: 'center',
        },
        {
          text: 'Fingerprint',
          value: 'data',
          align: 'center',
        },
        {
          text: 'Actions',
          value: 'actions',
          align: 'center',
        },
      ],
    };
  },

  computed: {
    getListPrivateKeys() {
      return this.$store.getters['privatekeys/list'];
    },

    getNumberPrivateKeys() {
      return this.$store.getters['privatekeys/getNumberPrivateKeys'];
    },
  },

  created() {
    this.dialog = !(localStorage.getItem('privatekeyPrivacyPolicy') === 'true');
  },

  methods: {
    convertToFingerprint(privateKey) {
      return parsePrivateKey(privateKey).fingerprint('md5');
    },

    accept() {
      localStorage.setItem('privatekeyPrivacyPolicy', this.privatekeyPrivacyPolicy);
      this.dialog = false;
    },
  },
};
</script>
