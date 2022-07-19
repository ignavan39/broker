import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { connectionService } from "../../api/Connection";
import { invitationService } from "../../api/Invitation";
import { errorState } from "../../state/Error.state";
import { userState } from "../../state/User.state";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: column;
`;

const Header = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  text-align: center;
  font-size: 2rem;
  color: white;
  height: 4rem;
  border-radius: 10px 10px 0 0;
  background-color: #4caf50;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  width: 23rem;
  justify-content: space-between;
  min-height: 400px;
  border-radius: 10px;
  border: 1px solid #4caf50;
`;

const FormButton = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  min-width: 100%;
  margin: 0.2rem 0;
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

export const Invitation = () => {
  const [err, setErr] = useRecoilState(errorState);

  const [connection, setConnection] = useState<{
    queueName: string;
    exchange: string;
    host: string;
    port: number | null;
    user: string;
    vhost: string;
    password: string;
  }>({
    queueName: "",
    exchange: "",
    host: "",
    port: null,
    user: "",
    vhost: "",
    password: "",
  });

  const navigate = useNavigate();

  const splits = window.location.href.split("/");
  const code = splits[splits.length - 1];

  useEffect(() => {
    (async () => {
      try {
        const response = await connectionService.connect();
        setConnection({
          queueName: response.queueName,
          exchange: response.exchange,
          host: response.host,
          port: response.port,
          user: response.user,
          vhost: response.vhost,
          password: response.password,
        })
        console.log(connection);
        connectionService.getData(code, connection.host, connection.port, connection.queueName);
      } catch (e) {
        const message = e instanceof Error ? e.message : "unknown error";
        setErr(message);
      }
    })();
  }, [setErr]);

  const invitationText = "You have been invited to " + code;

  return (
    <Container>
      <Form>
        <Header>Invitation</Header>
        <div style={{ margin: "2rem" }}>{invitationText}</div>
        <FormButton>
          <Button onClick = {() => {
            invitationService.accept(code)
          }}>
            Accept
          </Button>
          <Button onClick = {() => {
            invitationService.reject(code)
          }}>
            Reject
          </Button>
        </FormButton>
      </Form>
    </Container>
  )
};
