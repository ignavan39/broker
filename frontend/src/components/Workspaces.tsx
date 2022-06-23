import { useEffect, useState } from "react";
import styled from "styled-components";
import { getWorkspaces } from "../api/Workspace";
import { WorkspaceList } from "../types/Worpkspace";
import { ErrorPopup } from "./ErrorPopup";
import { Navbar } from "./Navbar";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: row;
`;

export const Workspace = () => {
  const [workspaces, setWorkspaces] = useState<WorkspaceList>([]);
  const [err, setErr] = useState<string | null>(null);
  const [errorPopupState, setOpenPopupState] = useState<boolean>(false);

  useEffect(() => {
    (async () => {
      try {
        let response = await getWorkspaces();
        setWorkspaces(response)
      } catch (e) {
        const message = e instanceof Error ? e.message : "unknown error";
        setOpenPopupState(true);
        setErr(message);
      }
    })();
  }, []);
  return (
    <>
      <Navbar />
      <Container>
        {errorPopupState && err ? (
          <ErrorPopup err={err} setOpen={setOpenPopupState} />
        ) : (
          <></>
        )}
        {}
      </Container>
    </>
  );
};
