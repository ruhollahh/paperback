export function getCart() {
  return JSON.parse(localStorage.getItem("cart")) || [];
}

export function setCart(cart) {
  localStorage.setItem("cart", JSON.stringify(cart));
}

export async function addToCart(id) {
  const product = await getProductById(id);
  const results = app.store.cart.filter(
    (prodInCart) => prodInCart.product.id == id
  );
  if (results.length == 1) {
    app.store.cart = app.store.cart.map((p) =>
      p.product.id == id ? { ...p, quantity: p.quantity + 1 } : p
    );
  } else {
    app.store.cart = [...app.store.cart, { product, quantity: 1 }];
  }
}

export function removeFromCart(id) {
  app.store.cart = app.store.cart.filter(
    (prodInCart) => prodInCart.product.id != id
  );
}
