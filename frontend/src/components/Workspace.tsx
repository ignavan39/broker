import styled from "styled-components";
import { WorkspaceListItem } from "../types/Workspace";

const Container = styled.div`
  min-height: 5rem;
  width: 28rem;
  border-radius: 15px;
  display: flex;
  background-color: #dedede;
  flex-direction: column;
  justify-content: center;
  padding:0 0.5rem;
  align-items: space-between;
  margin: 1.5rem 0;
`;
export const WorkspaceItem = ({
  id,
  name,
  isPrivate,
  createdAt,
}: WorkspaceListItem) => {
  return (
    <>
      <Container style={isPrivate ? {backgroundColor:'#96b3e3'} : {backgroundColor:'#dedede'}}>
        {name}
        </Container>

    </>
  );
};
