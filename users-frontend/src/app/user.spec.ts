import { UserStatus, userStatusToString } from "./user";

describe('User', () => {
    it('should return Active', () => {
      expect(userStatusToString(UserStatus.Active)).toBe("Active");
    });

    it('should return Inactive', () => {
        expect(userStatusToString(UserStatus.Inactive)).toBe("Inactive");
    });

    it('should return Terminated', () => {
    expect(userStatusToString(UserStatus.Terminated)).toBe("Terminated");
    });
  });