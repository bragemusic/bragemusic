import { BrowserRouter, Route, Routes } from "react-router-dom";

import { useEffect, useRef, useState } from "react";
import { useMediaQuery } from "react-responsive";

import { Player } from "./components/player";
import { Menu } from "./components/Menu";
import { MobileMenu } from "./components/MobileMenu";
import { MobilePlayer } from "./components/MobilePlayer";
import ArtistsPage from "./pages/artists";
import ArtistPage from "./pages/ArtistPage";
import AlbumPage from "./pages/AlbumPage";
import Toast from "./components/toast";
import AlbumsPage from "./pages/AlbumsPage";
import PlaylistsPage from "./pages/PlaylistsPage";
import PlaylistPage from "./pages/PlaylistPage";
import { Login } from "./components/Login";
import ImportPage from "./pages/ImportPage";
import ServerAdminPage from "./pages/ServerAdminPage";
import { SessionProvider, useSession } from "./session/SessionContext";
import LikedTracksPage from "./pages/LikedTracksPage";
import { usePlayerState } from "./store/playerStore";
import { useDeviceStore } from "./store/deviceStore";
import { useApi, usePlayerApi } from "./api/ApiContext";
import MediaPage from "./pages/MediaPage";
import HomePage from "./pages/HomePage";
import ServerPage from "./pages/ServerPage";
import TracksPage from "./pages/TracksPage";
import UserSettingsPage from "./pages/UserSettingsPage";

import { mqMobile } from "@/config/config";
import AboutPage from "@/pages/about";
import BlogPage from "@/pages/blog";
import { ThemeProvider } from "next-themes";
import { Spinner } from "@heroui/react";
import { LocalPlayer } from "./components/LocalPlayer";
import { MediaSession } from "./components/MediaSession";

function App() {
  return (
    <ThemeProvider>
      <SessionProvider>
        <BrowserRouter>
          <>
            <LocalPlayer/>
            <MediaSession/>
            <AppContent />
          </>
        </BrowserRouter>
      </SessionProvider>
    </ThemeProvider>
  );
}

function AppContent() {
  const [showLogin, setShowLogin] = useState(false);
  const [showSpinner, setShowSpinner] = useState(true);

  const api = useApi();
  const session = useSession();
  const scrollRef = useRef<HTMLDivElement>(null);
  const isMobile = useMediaQuery({ query: mqMobile });

  useEffect(() => {
    usePlayerState.getState().init(api);
    useDeviceStore.getState().init(api);
    api.getUser()
  }, []);

  useEffect(() => {
    setShowLogin(session.user === null);
    setShowSpinner(session.user === undefined)
    useDeviceStore.getState().loadDevices(api);
  }, [session.user]);

  let content;

  if (showSpinner) {
    content = (<div className="flex justify-center items-center w-full h-full"><Spinner size="xl"/></div>);
  } else if (showLogin) {
    content = <Login />;
  } else {
    content = (
        <>
          <Toast />

          <div className="flex flex-1 min-h-0">
            {!isMobile && <Menu />}

            <div className="flex relative flex-col flex-1 min-h-0">
              {isMobile && <MobilePlayer />}
              {isMobile && <MobileMenu />}

              <div
                ref={scrollRef}
                className="overflow-y-auto flex-1 min-h-0 bg-background"
              >
                <Routes>
                  <Route element={<HomePage  scrollRef={scrollRef}/>} path="/" />
                  <Route element={<ArtistsPage />} path="/artists" />
                  <Route element={<ArtistPage scrollRef={scrollRef} />} path="/artists/:artistID"/>
                  <Route element={<AlbumsPage />} path="/albums" />
                  <Route
                    element={<AlbumPage scrollRef={scrollRef} />}
                    path="/albums/:albumID"
                  />
                  <Route element={<ImportPage />} path="/import" />
                  <Route
                    element={<LikedTracksPage scrollRef={scrollRef} />}
                    path="/liked-tracks"
                  />
                  <Route element={<PlaylistsPage />} path="/playlists" />
                  <Route
                    element={<PlaylistPage scrollRef={scrollRef} />}
                    path="/playlists/:playlistID"
                  />
                  <Route element={<ServerAdminPage />} path="/server-admin" />
                  <Route element={<TracksPage />} path="/tracks" />
                  <Route element={<UserSettingsPage />} path="/settings" />
                  {isMobile && <Route element={<MediaPage />} path="/media" />}
                  {isMobile && <Route element={<ServerPage />} path="/server" />}
                  <Route element={<BlogPage />} path="/blog" />
                  <Route element={<AboutPage />} path="/about" />
                  <Route path="*" element={<div>Attans</div>} />
                </Routes>
              </div>

              {!isMobile && <Player />}
            </div>
          </div>
        </>
    );
  }

  return (
    <div className="flex overflow-hidden flex-col select-none h-dvh bg-background" >
      {content}
    </div>
  );
}

export default App;
