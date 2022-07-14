import styled from "styled-components";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: column;
`;

const Button = styled.button`
  background-color: #4caf50;
  border: none;
  color: white;
  padding: 1rem 2rem;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  min-width: 9rem;
  margin: 0.2rem 1rem;
  min-height: 1rem;
  border-radius: 10px;
  &:hover {
    background-color: #4aaf90;
  }
`;

const FormButton = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  min-width: 100%;
  margin: 0.2rem 0;
`;

export const Invitation = () => {
  return <Container></Container>
}