import { getProducts } from "./services/cart.js";
import store from "./services/store.js";

window.app = { store };

window.addEventListener("DOMContentLoaded", () => {
  app.store.cart = getProducts();
});
