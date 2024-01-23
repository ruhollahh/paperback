import { helpme } from "../helper/helpers";

export default class CartComp extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    console.log(helpme());
  }
}

customElements.define("cart-comp", CartComp);
