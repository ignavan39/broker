import { useNavigate } from "react-router";
import { useRecoilState } from "recoil";
import styled from "styled-components";
import { userState } from "../state/User.state";

const Container = styled.div`
  position: absolute;
  height: 8vh;
  width: 100vw;
  display: flex;
  justify-content: right;
  align-items: center;
  background-color: #e3e3e3;
  margin:0;
`;
const UserInfo = styled.button`
  display: flex;
  align-items: center;
  justify-content: start;
  background-color: white;
  min-width: 15rem;
  border-radius: 10px;
  height: 85%;
  border: none;
  font-size: 1rem;
  &:hover {
    border: 2px solid #000;
  }
`;
const Avatar = styled.div`
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  display: inline-block;
  vertical-align: middle;
  margin: 0.1rem 0.5rem;
`;

const Logout = styled.button`
  background-color: #8a8a8a;
  border: none;
  color: white;
  padding: 1rem 2rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: start;
  font-size: 18px;
  min-width: 2rem;
  margin: 0.2rem 2rem;
  height: 70%;
  border-radius: 10px;
  &:hover {
    background-color: #ff2626;
  }
`;

export const Navbar = () => {
  const [user, setUser] = useRecoilState(userState);
  const navigate = useNavigate();

  const logout = () => {
    setUser({
      ...user,
      auth: {
        accessToken: "",
      },
    });
    localStorage.removeItem("user");
    navigate("/auth");
  };
  return (
    <Container>
      <UserInfo>
        <Avatar>
          <img
            src="https://vk.com/images/camera_c.gif"
            width="50"
            height="50"
          />
        </Avatar>
        <div>{user.firstName + " " + user.lastName}</div>
      </UserInfo>
      <Logout onClick={logout}>Logout</Logout>
    </Container>
  );
};
