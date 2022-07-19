import { atom } from "recoil";
import { Invitation } from "../types/Invitation";

const getDefaultInvitation = () => {
    const cache = localStorage.getItem("invitation");
    if (!cache) {
        return {
            id: "",
            workspace: {
                id: "",
                name: "",
                isPrivate: null,
                createdAt: null,
            },
            senderId: "",
            status: null, 
            systemStatus: null,
            createdAt: "",
        }
    } else {
        const invitation = JSON.parse(cache) as Invitation;
        return invitation;
    }
}

export const invitationState = atom<Invitation>({
    key: "InvitationState",
    default: getDefaultInvitation()
})