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
import { userLoggedOn } from "./state/User.state";
import { Invitation } from "./components/invitation/Invitation";
import { invitationState } from "./state/Invitation.state"

const BaseRouter = () => {
  const loggedOn = useRecoilValue(userLoggedOn);
  const [err, setErr] = useRecoilState(errorState);
  const [isInvitation, invState] = useRecoilState(invitationState)
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
          <Route path="/invitations/*" element={<Invitation />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Registration />} />
          <Route
            path="*"
            element={loggedOn ? <Home /> : <Navigate to="/" />}
          />
          <Route
            path="/"
            element={loggedOn ? <Workspaces /> : <Navigate to="/login" />}
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
