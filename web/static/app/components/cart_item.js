// import { removeFromCart } from "../services/cart.js";

export default class CartItem extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {}
}

customElements.define("cart-item", CartItem);
