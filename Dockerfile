FROM alpine
ADD cart-api-api /cart-api-api
ENTRYPOINT [ "/cart-api-api" ]
