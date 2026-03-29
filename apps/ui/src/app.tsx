import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import { MetaProvider, Title, Link } from "@solidjs/meta";
import { ColorModeProvider } from "~/lib/color-mode";
import "./app.css";

export default function App() {
  return (
    <MetaProvider>
      <Title>PocketPanel - Admin</Title>
      <Link rel="icon" href="/favicon.svg" />
      <ColorModeProvider>
        <Router root={(props) => <Suspense>{props.children}</Suspense>}>
          <FileRoutes />
        </Router>
      </ColorModeProvider>
    </MetaProvider>
  );
}
