import styled from "styled-components";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50rem;
  width: 20rem;
  margin-left:1rem;
  margin-top: 3vh;
  border-radius: 15px;
  flex-direction: column;
`;

export const Settings = () => {
  return <Container></Container>;
};
