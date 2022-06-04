import styled from "styled-components";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50rem;
  width: 90rem;
  border-radius: 10px;
  margin-top: 3vh;
  margin-right: 5rem;
  margin-left: 0.5rem;
  flex-direction: row;
`;

const ChatContainer = styled.div`
  height: 50rem;
  width: 90rem;
  border-radius: 15px 0 0 15px;
  display: flex;
  background-color: #dedede;
`;

const Dialogues = styled.div`
  display: flex;
  justify-content: center;
  align-items: space-between;
  min-height: 50rem;
  min-width: 20rem;
  border-radius: 0 15px 15px 0;
  background-color: #dedede;
  border-left: 2px solid #cecece;
  flex-direction: column;
`;

export const Chat = () => {
  return (
    <>
      <Container>
        <ChatContainer></ChatContainer>
        <Dialogues></Dialogues>
      </Container>
    </>
  );
};
