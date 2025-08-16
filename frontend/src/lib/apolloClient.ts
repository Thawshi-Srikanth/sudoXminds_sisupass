import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";


const GRAPHQL_URI =
  process.env.NEXT_PUBLIC_GRAPHQL_URI || "http://localhost:8000/graphql/";

const client = new ApolloClient({
  link: new HttpLink({
    uri: GRAPHQL_URI,
    credentials: "include", // send cookies
  }),
  cache: new InMemoryCache(),
});

export default client;
