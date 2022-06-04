import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { userState } from "../state/User.state";
import { User } from "../types/User";
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
  width: 400px;
  justify-content: space-between;
  height: 300px;
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
  font-size: 20px;
  color: white;
  height: 4rem;
  border-radius: 10px 10px 0 0;
  background-color: #4caf50;
`;

const Button = styled.button`
  background-color: #4caf50;
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  min-width: auto;
  margin: 0.5rem 1rem;
  min-height: 2.5rem;
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
  min-height: 3rem;
  font-size: 16px;
  text-align: center;
  margin: 0.2rem 0;
`;

export const Auth = () => {
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState({
    password: user.password ?? "",
    email: user.email ?? "",
  });
  const navigate = useNavigate();

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    setState({
      ...state,
      [e.target.name]: e.target.value,
    });
  };

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const newUser: User = {
      ...user,
      email: state.email,
      password: state.password,
      auth: {
        accessToken: Math.random().toString(),
        expiresAt: new Date(),
      },
    };
    setUser(newUser);
    localStorage.setItem("user", JSON.stringify(newUser));
    navigate("/");
  };
  return (
    <Container>
      <Form onSubmit={onSubmit}>
        <Header>Login</Header>
        <FormInput>
          <Input
            type={"email"}
            placeholder="email"
            onInput={handleInput}
            value={state.email}
            name="email"
          />
          <Input
            type={"password"}
            minLength={5}
            placeholder="password"
            onInput={handleInput}
            value={state.password}
            name="password"
          />
        </FormInput>
        <Button>Submit</Button>
      </Form>
    </Container>
  );
};
