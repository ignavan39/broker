import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { sign } from "../../api";
import { userState } from "../../state/User.state";
import { User } from "../../types/User";
import { ErrorPopup } from "../ErrorPopup";
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

export const Login = () => {
  const [err, setErr] = useState<string | null>(null);
  const [errorPopupState, setOpenPopupState] = useState<boolean>(false);
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState({
    password: user.user.password ?? "",
    email: user.user.email ?? "",
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
        if (user.user.password.length && user.user.email.length) {
          const apiResponse = await sign({
            password: user.user.password,
            email: user.user.email,
            operation: "sign_in",
          });
          const updatedUser: User = {
            ...apiResponse,
            user: {
              ...apiResponse.user,
              password: state.password,
            },
          };
          setUser(updatedUser);
          localStorage.setItem("user", JSON.stringify(updatedUser));
          navigate("/");
        }
      } catch (e) {
        const message = e instanceof Error ? e.message : "unknown error";
        setOpenPopupState(true);
        setErr(message);
      }
    })();
  }, []);

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const apiResponse = await sign({
        password: state.password,
        email: state.email,
        operation: "sign_in",
      });
      const updatedUser: User = {
        ...apiResponse,
        user: {
          ...apiResponse.user,
          password: state.password,
        },
      };
      setUser(updatedUser);
      localStorage.setItem("user", JSON.stringify(apiResponse));
      navigate("/");
    } catch (e) {
      const message = e instanceof Error ? e.message : "unknown error";
      setOpenPopupState(true);
      setErr(message);
    }
  };
  return (
    <Container>
      {errorPopupState && err ? (
        <ErrorPopup err={err} setOpen={setOpenPopupState} />
      ) : (
        <></>
      )}
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