ARG CILINT_VERSION

FROM golangci/golangci-lint:${CILINT_VERSION}

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

CMD ["golangci-lint", "run", "-v"]
