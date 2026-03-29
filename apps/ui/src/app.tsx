import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import { isServer } from "solid-js/web";
import { MetaProvider, Title, Link } from "@solidjs/meta";
import { ColorModeProvider, ColorModeScript, cookieStorageManagerSSR } from "@kobalte/core";
import { getCookie } from "@solidjs/start/http";
import { ModeToggle } from "~/components/ModeToggle";
import "./app.css";

function getServerCookies() {
  "use server";
  const colorMode = getCookie("kb-color-mode");
  return colorMode ? `kb-color-mode=${colorMode}` : "";
}

export default function App() {
  const storageManager = cookieStorageManagerSSR(isServer ? getServerCookies() : document.cookie);
  return (
    <MetaProvider>
      <Title>PocketPanel - Admin</Title>
      <Link rel="icon" href="/favicon.svg" />
      <ColorModeScript storageType={storageManager.type} />
      <ColorModeProvider storageManager={storageManager}>
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
      </ColorModeProvider>
    </MetaProvider>
  );
}
