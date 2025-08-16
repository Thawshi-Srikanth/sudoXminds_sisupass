import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";

const client = new ApolloClient({
  link: new HttpLink({
    uri: "http://192.168.142.225:8000/graphql/",
    credentials: "include", // send cookies
  }),
  cache: new InMemoryCache(),
});

export default client;
