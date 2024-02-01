# Dockerfile to run prettier on markdown in the codebase. Used by lint-markdown
# and lint-markdown-fix rules in Makefile.

FROM node:alpine

WORKDIR /usr/spin-operator

RUN npm install prettier -g

ENV PRETTIER_MODE=check

CMD sh -c "npx prettier --${PRETTIER_MODE} '**/*.md'"
