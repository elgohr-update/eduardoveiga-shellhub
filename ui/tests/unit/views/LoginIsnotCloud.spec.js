import Vuex from 'vuex';
import { mount, createLocalVue, config } from '@vue/test-utils';
import Vuetify from 'vuetify';
import Login from '@/views/Login';
import router from '@/router/index';

config.mocks = {
  $env: {
    isCloud: false,
  },
};

describe('Login', () => {
  const localVue = createLocalVue();
  localVue.use(Vuex);
  const vuetify = new Vuetify();

  let wrapper;

  const auth = {
    status: true,
    token: 'xxxxxxxx',
    user: 'user',
    tenant: 'xxxxxxxx',
  };

  const store = new Vuex.Store({
    state: {
      auth,
    },
    getters: {
    },
    actions: {
      'auth/logout': () => {},
      'auth/login': () => {},
      'auth/loginToken': () => {},
      'notifications/fetch': () => {},
      'snackbar/showSnackbarErrorIncorrect': () => {},
      'snackbar/showSnackbarErrorDefault': () => {},
    },
  });

  describe('Is cloud', () => {
    beforeEach(() => {
      wrapper = mount(Login, {
        store,
        localVue,
        router,
        mocks: {
          $route: {
            query: {},
          },
        },
        vuetify,
      });
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue instance', () => {
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Compare data with default value', () => {
      expect(wrapper.vm.username).toEqual('');
      expect(wrapper.vm.password).toEqual('');
      expect(wrapper.vm.error).toEqual(false);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with components', () => {
      expect(wrapper.find('[data-test="accountCreated-component"]').exists()).toEqual(false);
    });
    it('Renders the template with data', () => {
      expect(wrapper.find('[data-test="username-text"]').element.value).toEqual('');
      expect(wrapper.find('[data-test="password-text"]').element.value).toEqual('');
      expect(wrapper.find('[data-test="login-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="forgotPassword-card"]').exists()).toBe(false);
      expect(wrapper.find('[data-test="isCloud-card"]').exists()).toBe(false);
    });
  });
});
