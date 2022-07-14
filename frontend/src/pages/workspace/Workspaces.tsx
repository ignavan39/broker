import { useEffect, useState } from "react";
import styled from "styled-components";
import { WorkspaceList } from "../../types/Workspace";
import { ErrorPopup } from "../../components/ErrorPopup";
import { Navbar } from "../../components/Navbar";
import { WorkspaceItem } from "../../components/workspace/Workspace";
import { workspaceService } from "../../api/Workspace";
import { useRecoilState } from "recoil";
import { errorState } from "../../state/Error.state";

const Container = styled.div`
  display: flex;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: column;
  background-color: white;
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
  height: 1rem;
  border-radius: 10px;
  &:hover {
    background-color: #ff2626;
  }
`;

const WorkspaceListContainer = styled.div`
  display: flex;
  height: 60vh;
  width: 25vw;
  flex-direction: column;
  overflow-y: scroll;
  position: absolute;
  padding: 5px;
  top: 2rem;
  margin: 15rem 0;
  background-color: white;
  border-radius: 10px;
`;

const CreateWorkspaceContainer = styled.form`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 10vh;
  width: 30rem;
  margin: 1rem 0;
  position: absolute;
  top: 2rem;
  margin: 5rem;
  background-color: #dedede;

  border-radius: 10px;
`;

export const Workspaces = () => {
  const [workspaces, setWorkspaces] = useState<WorkspaceList>(
    [] as WorkspaceList
  );
  const [err, setErr] = useRecoilState(errorState);
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
      let response = await workspaceService.create(state.name, state.isPrivate);
      setWorkspaces([...workspaces, response].sort((p,n)=> p.createdAt < n.createdAt ? 1 : 0));
    } catch (e) {
      const message = e instanceof Error ? e.message : "unknown error";
      setErr(message);
    }
  };

  useEffect(() => {
    (async () => {
      try {
        let response = await workspaceService.getAll();
        setWorkspaces(response.workspaces);
      } catch (e) {
        const message = e instanceof Error ? e.message : "unknown error";
        setErr(message);
      }
    })();
  }, [setErr]);

  return (
    <>
      <Navbar />
      <Container>
        <CreateWorkspaceContainer onSubmit={onSubmit}>
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
        </CreateWorkspaceContainer>
        <WorkspaceListContainer>
          {workspaces.map((w) => (
            <WorkspaceItem
              id={w.id}
              createdAt={w.createdAt}
              name={w.name}
              isPrivate={w.isPrivate}
            />
          ))}
        </WorkspaceListContainer>
      </Container>
    </>
  );
};
