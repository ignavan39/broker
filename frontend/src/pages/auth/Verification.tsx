import { useNavigate } from "react-router";
import { userState } from "../../state/User.state";
import styled from "styled-components";
import { useRecoilState } from "recoil";
import { errorState } from "../../state/Error.state";
import { useState } from "react";
import { User } from "../../types/User";
import { authorizationService } from "../../api";
import axios, { AxiosError } from "axios";

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

const Input = styled.input`
  min-height: 3rem;
  border: 1px solid #bdbdbd;
  border-radius: 10px;
  padding: 0 10px;
  min-width: auto;
  font-size: 16px;
  text-align: center;
  margin: 0.2rem 0;
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
      if (!state.code) {
        throw new Error("Invalid Invitation Code");
      }

      const apiResponse = await authorizationService.signUp({
        password: user.user.password,
        email: user.user.email,
        code: Number.parseInt(state.code),
        lastName: user.user.lastName,
        firstName: user.user.firstName,
        nickname: user.user.nickname,
      });
      const updatedUser: User = {
        ...apiResponse,
        user: {
          ...apiResponse.user,
          password: user.user.password,
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

