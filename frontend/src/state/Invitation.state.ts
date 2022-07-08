import { atom } from "recoil";
import { Invitation } from "../types/Invitation";

const getDefaultInvitation = () => {
    const cache = localStorage.getItem("invitation");
    if (!cache) {
        return {
            id: "",
            workspaceId: "",
            senderId: "",
            status: "", 
            systemStatus: "",
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