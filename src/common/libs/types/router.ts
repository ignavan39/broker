import { Router } from "express";

export interface InternalRouter {
    getRouter(): Router
}