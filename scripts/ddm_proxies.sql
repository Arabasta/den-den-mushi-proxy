CREATE TABLE `ddm_proxies` (
                               `HostName` varchar(191) NOT NULL,
                               `ProxyType` varchar(191) DEFAULT NULL,
                               `Url` varchar(191) NOT NULL,
                               `IpAddress` varchar(191) DEFAULT NULL,
                               `Environment` varchar(191) NOT NULL DEFAULT 'development',
                               `Country` varchar(191) NOT NULL DEFAULT 'unknown',
                               `LastHeartbeatAt` timestamp NULL DEFAULT NULL,
                               `DrainingStartAt` timestamp NULL DEFAULT NULL,
                               `DrainingEndAt` timestamp NULL DEFAULT NULL,
                               `DeploymentColor` varchar(191) NOT NULL DEFAULT 'blue',
                               `MaxSessions` bigint DEFAULT '100',
                               `ActiveSessions` bigint DEFAULT '0',
                               PRIMARY KEY (`HostName`),
                               UNIQUE KEY `uni_ddm_proxies_url` (`Url`),
                               UNIQUE KEY `uni_ddm_proxies_ip_address` (`IpAddress`)
)

    INSERT INTO ddm_proxies
    (HostName, ProxyType, Url, IpAddress, Environment, Country)
    VALUES (
    'Kes-MacBook-Pro.local', 'OS', 'http://localhost:45007',
            '127.0.0.1', 'Development', 'SG');
