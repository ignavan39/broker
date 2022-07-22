import { Button, Input } from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { authorizationService } from "../../api";
import { errorState } from "../../state/Error.state";
import { userState } from "../../state/User.state";
import { User } from "../../types/User";
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
`;

export const Login = () => {
  const [err, setErr] = useRecoilState(errorState);
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState({
    password: user.profile.password ?? "",
    email: user.profile.email ?? "",
  });
  const navigate = useNavigate();

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    setState({
      ...state,
      [e.target.name]: e.target.value,
    });
  };

  useEffect(() => {
    (async () => {
      try {
        if (user.profile.password.length && user.profile.email.length) {
          const apiResponse = await authorizationService.signIn({
            password: user.profile.password,
            email: user.profile.email,
          });
          const updatedUser: User = {
            ...apiResponse,
            profile: {
              ...apiResponse.profile,
              password: state.password,
            },
          };
          setUser(updatedUser);
          localStorage.setItem("user", JSON.stringify(updatedUser));
          navigate("/workspaces");
        }
      } catch (e) {
        const message = e instanceof Error ? e.message : "unknown error";
        setErr(message);
      }
    })();
  }, []);

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const apiResponse = await authorizationService.signIn({
        password: state.password,
        email: state.email,
      });
      const updatedUser: User = {
        ...apiResponse,
        profile: {
          ...apiResponse.profile,
          password: state.password,
        },
      };
      setUser(updatedUser);
      localStorage.setItem("user", JSON.stringify(apiResponse));

      navigate("/workspaces");
    } catch (e) {
      const message = e instanceof Error ? e.message : "unknown error";
      setErr(message);
    }
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
            placeholder="password"
            onInput={handleInput}
            value={state.password}
            name="password"
          />
        </FormInput>
        <FormButton>
          <Button>Submit</Button>
          <Button
            style={{ backgroundColor: "#ff9900" }}
            onClick={() => {
              navigate("/register");
            }}
          >
            Registration
          </Button>
        </FormButton>
      </Form>
    </Container>
  );
};
