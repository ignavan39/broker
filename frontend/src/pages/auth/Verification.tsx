import { useNavigate } from "react-router";
import { userState } from "../../state/User.state";
import styled from "styled-components";
import { useRecoilState } from "recoil";
import { errorState } from "../../state/Error.state";
import { useState } from "react";
import { User } from "../../types/User";
import { authorizationService } from "../../api";
import axios, { AxiosError } from "axios";
import { connectionService } from "../../api/Connection";
import { Button, Input } from "@mui/material";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: column;
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

const FormInput = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin: 0 1rem;
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

const FormButton = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  min-width: 100%;
  margin: 0.2rem 0;
`

export const Verification = () => {
  const [err, setErr] = useRecoilState(errorState);
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState<{
    code: string;
  }>({
    code: "",
  });
  const navigate = useNavigate();

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const apiResponse = await authorizationService.signUp({
        password: user.profile.password,
        email: user.profile.email,
        code: Number.parseInt(state.code),
        lastName: user.profile.lastName,
        firstName: user.profile.firstName,
        nickname: user.profile.nickname,
      });
      const updatedUser: User = {
        ...apiResponse,
        profile: {
          ...apiResponse.profile,
          password: user.profile.password,
        },
      };
      setUser(updatedUser);
      localStorage.setItem("user", JSON.stringify(apiResponse));

      navigate("/workspaces");
    } catch (e) {
      const message = e instanceof Error ? e.message : "unknown error";
      setErr(message);
    }
  }

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    setState({
      ...state,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <Container>
      <Form onSubmit={onSubmit}>
        <Header>Code verification</Header>
        <FormInput>
          <Input
            type={"code"}
            placeholder="code"
            onInput={handleInput}
            value={state.code}
            name="code"
          />
        </FormInput>
        <FormButton>
          <Button>Submit</Button>
        </FormButton>
      </Form>
    </Container>
  )
}

