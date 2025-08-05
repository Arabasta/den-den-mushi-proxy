export * from './healthcheck.service';
import { HealthcheckService } from './healthcheck.service';
export * from './makeChange.service';
import { MakeChangeService } from './makeChange.service';
export * from './pTYToken.service';
import { PTYTokenService } from './pTYToken.service';
export * from './whitelistBlacklist.service';
import { WhitelistBlacklistService } from './whitelistBlacklist.service';
export const APIS = [HealthcheckService, MakeChangeService, PTYTokenService, WhitelistBlacklistService];
