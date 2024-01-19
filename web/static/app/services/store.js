const store = {
  cart: [],
};

const proxiedStore = new Proxy(store, {
  set(target, property, value) {
    target[property] = value;
    if (property == "cart") {
      window.dispatchEvent(new Event("appcartchange"));
    }
    return true;
  },
  get(target, property) {
    return target[property];
  },
});

export default proxiedStore;
