import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import { MetaProvider, Title, Link } from "@solidjs/meta";
import { ModeToggle } from "~/components/ModeToggle";
import "./app.css";

export default function App() {
  return (
    <MetaProvider>
      <Title>PocketPanel - Admin</Title>
      <Link rel="icon" href="/favicon.svg" />
      <Router
        root={props => (
          <>
            <ModeToggle />
            <Suspense>{props.children}</Suspense>
          </>
        )}
      >
        <FileRoutes />
      </Router>
    </MetaProvider>
  );
}
