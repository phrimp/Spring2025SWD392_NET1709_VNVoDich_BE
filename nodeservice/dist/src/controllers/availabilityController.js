"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __rest = (this && this.__rest) || function (s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getCourseAvailability = exports.updateAvailability = exports.getTutorAvailability = void 0;
const availabilityService_1 = require("../services/availabilityService");
const availabilityMessage_1 = require("../message/availabilityMessage");
const getTutorAvailability = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { userId } = req.body;
        const availabilityData = yield (0, availabilityService_1.getTutorAvailabilityService)(userId);
        // if (!availabilityData) {
        //   res.json({ message: MESSAGES.tutorNotFound, data: null });
        //   return;
        // }
        res.json({
            message: availabilityMessage_1.MESSAGES.availabilityRetrieved,
            data: availabilityData,
        });
    }
    catch (error) {
        res.status(500).json({
            message: error.message || availabilityMessage_1.MESSAGES.errorRetrievingAvailability,
            error,
        });
    }
});
exports.getTutorAvailability = getTutorAvailability;
const updateAvailability = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const _a = req.body, { userId } = _a, data = __rest(_a, ["userId"]);
        const updatedAvailability = yield (0, availabilityService_1.updateAvailabilityService)(userId, data);
        res.json({
            message: availabilityMessage_1.MESSAGES.availabilityUpdated,
            data: updatedAvailability,
        });
    }
    catch (error) {
        res.status(500).json({
            message: error.message || availabilityMessage_1.MESSAGES.errorUpdatingAvailability,
            error,
        });
    }
});
exports.updateAvailability = updateAvailability;
const getCourseAvailability = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { courseId } = req.params;
        const { type } = req.query;
        const availableDates = yield (0, availabilityService_1.getCourseAvailabilityService)(Number(courseId), type);
        res.json({
            message: availabilityMessage_1.MESSAGES.courseAvailabilityRetrieved,
            data: availableDates,
        });
    }
    catch (error) {
        res.status(500).json({
            message: error.message || availabilityMessage_1.MESSAGES.errorRetrievingAvailability,
            error,
        });
    }
});
exports.getCourseAvailability = getCourseAvailability;
