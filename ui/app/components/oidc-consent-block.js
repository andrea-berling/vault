/**
 * @module OidcConsentBlock
 * OidcConsentBlock components are used to...
 *
 * @example
 * ```js
 * <OidcConsentBlock @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {string} redirect - redirect is the URL where successful consent will redirect to
 * @param {string} code - code is the string required to pass back to redirect on successful OIDC auth
 * @param {string} [state] - state is a string which is required to return on redirect if provided, but optional generally
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class OidcConsentBlockComponent extends Component {
  @tracked didCancel = false;

  get win() {
    return this.window || window;
  }

  buildUrl(urlString, params) {
    try {
      let url = new URL(urlString);
      Object.keys(params).forEach(key => {
        if (params[key]) {
          url.searchParams.append(key, params[key]);
        }
      });
      return url;
    } catch (e) {
      console.debug('DEBUG: parsing url failed for', urlString);
      throw new Error('Invalid URL');
    }
  }

  @action
  handleSubmit(evt) {
    evt.preventDefault();
    let { redirect, ...params } = this.args;
    let redirectUrl = this.buildUrl(redirect, params);
    if (this.args.onSuccess) {
      // Used for testing, but also makes redirect override available
      this.args.onSuccess(redirect, params);
    } else {
      this.win.location.replace(redirectUrl);
    }
  }

  @action
  handleCancel(evt) {
    evt.preventDefault();
    this.didCancel = true;
  }
}