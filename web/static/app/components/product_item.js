// import { addToCart } from "../services/cart.js";

export default class ProductItem extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {}
}

customElements.define("product-item", ProductItem);
