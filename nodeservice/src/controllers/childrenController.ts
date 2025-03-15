import { Request, Response } from "express";
import {
  getChildren,
  getChild,
  createChild,
  updateChild,
  deleteChild,
} from "../services/childService";
import { childMessages } from "../message/childMessages";

export const getChildrenHandler = async (req: Request, res: Response) => {
  try {
    const { userId } = req.body;
    if (!userId) {
      res.status(400).json({ message: childMessages.PARENT_ID_REQUIRED });
      return;
    }

    const children = await getChildren(Number(userId));
    res.json({ message: childMessages.CHILDREN_RETRIEVED, data: children });
  } catch (error) {
    res.status(500).json({ message: childMessages.ERROR_RETRIEVING_CHILDREN, error });
  }
};

export const getChildHandler = async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    if (!id || isNaN(Number(id))) {
      res.status(400).json({ message: childMessages.INVALID_CHILD_ID });
      return;
    }

    const child = await getChild(Number(id));
    if (!child) {
      res.status(404).json({ message: childMessages.CHILD_NOT_FOUND });
      return;
    }

    res.json({ message: childMessages.CHILD_RETRIEVED, data: child });
  } catch (error) {
    res.status(500).json({ message: childMessages.ERROR_RETRIEVING_CHILDREN, error });
  }
};

export const createChildHandler = async (req: Request, res: Response) => {
  try {
    const { full_name, username, learning_goals, password, userId, date_of_birth } = req.body;

    if (!full_name || !username || !learning_goals || !password || !userId || !date_of_birth) {
      res.status(400).json({ message: childMessages.ALL_FIELDS_REQUIRED });
      return;
    }

    const newChild = await createChild(
      full_name,
      username,
      learning_goals,
      password,
      Number(userId),
      date_of_birth
    );

    res.json({ message: childMessages.CHILD_CREATED, data: newChild });
  } catch (error) {
    res.status(500).json({ message: childMessages.ERROR_CREATING_CHILD, error });
  }
};

export const updateChildHandler = async (req: Request, res: Response) => {
  try {
    const { id } = req.params;
    const { full_name, learning_goals, password, date_of_birth } = req.body;

    if (!id || isNaN(Number(id))) {
      res.status(400).json({ message: childMessages.INVALID_CHILD_ID });
      return;
    }

    const updatedChild = await updateChild(
      Number(id),
      full_name,
      learning_goals,
      password,
      date_of_birth
    );

    res.json({ message: childMessages.CHILD_UPDATED, data: updatedChild });
  } catch (error) {
    res.status(500).json({ message: childMessages.ERROR_UPDATING_CHILD, error });
  }
};

export const deleteChildHandler = async (req: Request, res: Response) => {
  try {
    const { id } = req.params;

    if (!id || isNaN(Number(id))) {
      res.status(400).json({ message: childMessages.INVALID_CHILD_ID });
      return;
    }

    await deleteChild(Number(id));
    res.json({ message: childMessages.CHILD_DELETED });
  } catch (error) {
    res.status(500).json({ message: childMessages.ERROR_DELETING_CHILD, error });
  }
};
