import { useEffect, useState } from "react";
import styled from "styled-components";
import { createWorkspace, getWorkspaces } from "../api/Workspace";
import { WorkspaceList } from "../types/Worpkspace";
import { ErrorPopup } from "./ErrorPopup";
import { Navbar } from "./Navbar";
import { WorkspaceItem } from "./Worpkspace";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: row;
`;

const Button = styled.button`
  background-color: #8a8a8a;
  border: none;
  color: white;
  padding: 1rem 2rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: start;
  font-size: 18px;
  min-width: 2rem;
  margin: 0.2rem 2rem;
  height: 70%;
  border-radius: 10px;
  &:hover {
    background-color: #ff2626;
  }
`;

export const Workspaces = () => {
  const [workspaces, setWorkspaces] = useState<WorkspaceList>(
    [] as WorkspaceList
  );
  const [err, setErr] = useState<string | null>(null);
  const [errorPopupState, setOpenPopupState] = useState<boolean>(false);
  const [state, setState] = useState({ name: "", isPrivate: false });

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    setState({
      ...state,
      [e.target.name]: e.target.value,
    });
  };

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      console.log(state);
      let response = await createWorkspace(state.name, state.isPrivate);
      workspaces.push(response);
    } catch (e) {
      const message = e instanceof Error ? e.message : "unknown error";
      setOpenPopupState(true);
      setErr(message);
    }
  };

  useEffect(() => {
    (async () => {
      try {
        let response = await getWorkspaces();
        setWorkspaces(response.workspaces);
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
        <>
          {workspaces.map((w) => (
            <WorkspaceItem
              id={w.id}
              createdAt={w.createdAt}
              name={w.name}
              isPrivate={w.isPrivate}
            />
          ))}
          <form onSubmit={onSubmit}>
            <input
              type={"text"}
              placeholder="name"
              onInput={handleInput}
              value={state.name}
              name="name"
            ></input>
            <input
              type={"checkbox"}
              placeholder="isPrivate"
              name="isPrivate"
              onChange={() => {
                setState({ ...state, isPrivate: !!state.isPrivate });
              }}
              value={state.isPrivate.toString()}
            ></input>
            <Button>Create</Button>
          </form>
        </>
      </Container>
    </>
  );
};
