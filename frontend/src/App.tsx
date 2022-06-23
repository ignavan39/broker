import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { RecoilRoot, useRecoilValue } from "recoil";
import { Auth } from "./components/Auth";
import { Home } from "./components/Home";
import { Workspace } from "./components/Workspaces";
import { userIsLoggined } from "./state/User.state";

const BaseRouter = () => {
  const isLoggined = useRecoilValue(userIsLoggined);
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/auth" element={<Auth register={false} />} />
        <Route path="/register" element={<Auth register={true} />} />
        <Route path="*" element={isLoggined? <Home/> : <Navigate to='/' />} />
        <Route path="/" element={isLoggined? <Workspace/> : <Navigate to='/auth' />} />
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
