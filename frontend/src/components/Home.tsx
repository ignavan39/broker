import styled from "styled-components";
import { Chat } from "./Chat";
import { Navbar } from "./Navbar";
import { Settings } from "./Settings";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: row;
`;

export const Home = () => (
  <>
    {" "}
    <Navbar />
    <Container>
      <Settings />
      <Chat />
    </Container>
  </>
);
