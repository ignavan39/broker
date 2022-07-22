import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { userState } from "../../state/User.state";
import { User } from "../../types/User";
import { ErrorPopup } from "../../components/ErrorPopup";
import { authorizationService } from "../../api";
import { errorState } from "../../state/Error.state";
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
`;

export const Registration = () => {
  const [err, setErr] = useRecoilState(errorState);
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState<{
    password: string;
    email: string;
    firstName: string;
    lastName: string;
    nickname: string;
  }>({
    password: user.profile.password ?? "",
    email: user.profile.email ?? "",
    firstName: user.profile.firstName ?? "",
    lastName: user.profile.lastName ?? "",
    nickname: user.profile.nickname ?? "",
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
    try {
      if (!state.email) {
        throw new Error("Invalid Email");
      }

      await authorizationService.sendVerifyCode({
        email: state.email,
      });

      setUser({
        auth: {
          access: {
            token: "",
            expireAt: null,
          },
          refresh: {
            token: "",
            expireAt: null,
          },
        },
        profile : {
          ...state,
          avatarUrl: ""
        },
      })

      navigate("/verification")
    } catch (e) {
      const message = e instanceof Error ? e.message : "unknown error";
      setErr(message);
    }
  };
  return (
    <Container>
      <Form onSubmit={onSubmit}>
        <Header>Register</Header>
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
          <Input
            type={"text"}
            placeholder="Nickname"
            onInput={handleInput}
            value={state.nickname}
            name="nickname"
          />
          <Input
            type={"text"}
            placeholder="First Name"
            onInput={handleInput}
            value={state.firstName}
            name="firstName"
          />
          <Input
            type={"text"}
            placeholder="Last Name"
            onInput={handleInput}
            value={state.lastName}
            name="lastName"
          />
        </FormInput>
        <FormButton>
          <Button type="submit">Submit</Button>
        </FormButton>
      </Form>
    </Container>
  );
};
