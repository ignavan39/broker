import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { RecoilRoot, useRecoilValue } from "recoil";
import { Login } from "./components/auth/Login";
import { Registration } from "./components/auth/Registration";
import { Home } from "./components/Home";
import { Workspaces } from "./components/Workspaces";
import { userIsLoggined } from "./state/User.state";

const BaseRouter = () => {
  const isLoggined = useRecoilValue(userIsLoggined);
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Registration />} />
        <Route path="*" element={isLoggined ? <Home /> : <Navigate to="/" />} />
        <Route
          path="/"
          element={isLoggined ? <Workspaces /> : <Navigate to="/login" />}
        />
      </Routes>
    </BrowserRouter>
  );
};
const App = () => (
  <RecoilRoot>
    <BaseRouter />
  </RecoilRoot>
);

export default App;
