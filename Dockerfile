
FROM alpine:3.2
COPY ./authServer /
CMD ["/authServer"]
