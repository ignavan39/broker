import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { sign } from "../api";
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
`;

type AuthProp = {
  register: boolean;
};

export const Auth = (prop: AuthProp) => {
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState({
    password: user.password ?? "",
    email: user.email ?? "",
    firstName: !prop.register ? user.firstName ?? "" : "",
    lastName: !prop.register ? user.lastName ?? "" : "",
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
    let apiResponse = await sign(state, prop.register ? "signUp" : "signIn");
    const newUser: User = {
      ...user,
      ...state,
      auth: {
        accessToken: apiResponse.auth.accessToken,
        refreshToken: apiResponse.auth.refreshToken,
      },
    };
    setUser(newUser);
    localStorage.setItem("user", JSON.stringify(newUser));
    navigate("/");
  };
  return (
    <Container>
      <Form onSubmit={onSubmit}>
        {prop.register ? <Header>Register</Header> : <Header>Login</Header>}
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
          {prop.register ? (
            <>
              <Input
                type={"text"}
                minLength={1}
                placeholder="First Name"
                onInput={handleInput}
                value={state.firstName}
                name="firstName"
              />
              <Input
                type={"text"}
                minLength={1}
                placeholder="Last Name"
                onInput={handleInput}
                value={state.lastName}
                name="lastName"
              />
            </>
          ) : (
            <></>
          )}
        </FormInput>
        <FormButton>
          <Button>Submit</Button>
          {!prop.register ? (
            <Button
              style={{ backgroundColor: "#ff9900" }}
              onClick={() => {
                navigate("/register");
              }}
            >
              Registration
            </Button>
          ) : (
            <></>
          )}
        </FormButton>
      </Form>
    </Container>
  );
};
