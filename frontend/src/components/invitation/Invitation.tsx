import { useEffect } from "react";
import { useNavigate } from "react-router";
import styled from "styled-components";
import { invitationService } from "../../api/Invitation";

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  flex-direction: column;
`;

export const Invitation = () => {
  const navigate = useNavigate();
  useEffect(() => {
    (async () => {
      try {
        const splits = window.location.href.split("/");
        const code = splits[splits.length - 1];
        await invitationService.accept(code);
        navigate("/");
      } catch (err) {}
    })();
  }, []);
  return <Container></Container>;
};
