import { atom } from "recoil";
import { Invitation } from "../types/Invitation";

const getDefaultInvitation = () => {
    const cache = localStorage.getItem("invitation");
    if (!cache) {
        return null
    } else {
        const invitation = JSON.parse(cache) as Invitation;
        return invitation;
    }
}

export const invitationState = atom<Invitation | null>({
    key: "InvitationState",
    default: getDefaultInvitation()
})