FROM alpine:3.13

COPY cmd/bin/ /app/bin/
CMD [ "/app/bin/travel-article-headings" ]