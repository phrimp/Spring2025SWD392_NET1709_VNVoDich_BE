import { Request, Response } from "express";
import {
  getParents,
  getParentById,
  updateParentProfile,
} from "../services/parentService";
import { PARENT_MESSAGES } from "../message/parentMessages";

export const handleGetParents = async (req: Request, res: Response) => {
  try {
    const { page = 1, pageSize = 10 } = req.query;
    const pageNum = parseInt(page as string, 10);
    const pageSizeNum = parseInt(pageSize as string, 10);

    const result = await getParents(pageNum, pageSizeNum);
    res.json({
      message: PARENT_MESSAGES.RETRIEVED_SUCCESS,
      data: result.parents,
      pagination: result.pagination,
    });
  } catch (error) {
    console.error(PARENT_MESSAGES.ERROR_RETRIEVING, error);
    res.status(500).json({
      message: PARENT_MESSAGES.ERROR_RETRIEVING,
      error: (error as Error).message,
    });
  }
};

export const handleGetParentById = async (req: Request, res: Response) => {
  const { id } = req.params;

  try {
    const parent = await getParentById(Number(id));

    if (!parent) {
      res.status(404).json({ message: PARENT_MESSAGES.NOT_FOUND });
      return;
    }

    res.json({ message: PARENT_MESSAGES.RETRIEVED_ONE_SUCCESS, data: parent });
  } catch (error) {
    console.error(PARENT_MESSAGES.ERROR_RETRIEVING_ONE, error);
    res.status(500).json({
      message: PARENT_MESSAGES.ERROR_RETRIEVING_ONE,
      error: (error as Error).message,
    });
  }
};

export const handleUpdateParentProfile = async (
  req: Request,
  res: Response
) => {
  const { id } = req.params;

  try {
    const updatedParent = await updateParentProfile(Number(id), req.body);

    if (!updatedParent) {
      res.status(404).json({ message: PARENT_MESSAGES.NOT_FOUND });
      return;
    }

    res.json({
      message: PARENT_MESSAGES.UPDATED_SUCCESS,
      data: updatedParent,
    });
  } catch (error) {
    res.status(500).json({
      message: PARENT_MESSAGES.ERROR_UPDATING,
      error: (error as Error).message,
    });
  }
};
