import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { RecoilRoot, useRecoilState, useRecoilValue } from "recoil";
import { Home } from "./components/chat/Home";
import { ErrorPopup } from "./components/ErrorPopup";
import { Login } from "./pages/auth/Login";
import { Registration } from "./pages/auth/Registration";
import { Workspaces } from "./pages/workspace/Workspaces";
import { errorState } from "./state/Error.state";
import { userIsLoggined } from "./state/User.state";

const BaseRouter = () => {
  const isLoggined = useRecoilValue(userIsLoggined);
  const [err, setErr] = useRecoilState(errorState);
  return (
    <>
      {" "}
      {err ? (
        <ErrorPopup err={err} setOpen={setErr} />
      ) : (
        <></>
      )}
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Registration />} />
          <Route
            path="*"
            element={isLoggined ? <Home /> : <Navigate to="/" />}
          />
          <Route
            path="/"
            element={isLoggined ? <Workspaces /> : <Navigate to="/login" />}
          />
        </Routes>
      </BrowserRouter>
    </>
  );
};
const App = () => (
  <RecoilRoot>
    <BaseRouter />
  </RecoilRoot>
);

export default App;
