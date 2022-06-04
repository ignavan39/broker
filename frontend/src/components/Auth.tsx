import { useState } from "react";
import { useRecoilState } from "recoil";
import { userState } from "../state/User.state";

export const Auth = () => {
  const [user, setUser] = useRecoilState(userState);
  const [state, setState] = useState({
    password: "",
    email: user.email ?? "",
  });

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    setState({
      ...state,
      [e.target.name]: e.target.value,
    });
  };

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setUser({
      ...user,
      email: state.email,
    });
    localStorage.setItem("user", JSON.stringify(user));
  };
  return (
    <div>
      <form onSubmit={onSubmit}>
        <label>email</label>
        <input
          placeholder="email"
          onInput={handleInput}
          value={state.email}
          name="email"
        />
        <label>password</label>
        <input
          placeholder="password"
          onInput={handleInput}
          value={state.password}
          name="password"
        />
        <button />
      </form>
    </div>
  );
};
