/* tslint:disable */
/* eslint-disable */
/**
 * HumanLayer Daemon REST API
 * REST API for HumanLayer daemon operations, providing session management, approval workflows, and real-time event streaming capabilities.
 *
 * The version of the OpenAPI document: 1.0.0
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  HealthResponse,
} from '../models/index';
import {
    HealthResponseFromJSON,
    HealthResponseToJSON,
} from '../models/index';

/**
 * SystemApi - interface
 *
 * @export
 * @interface SystemApiInterface
 */
export interface SystemApiInterface {
    /**
     * Check if the daemon is running and healthy
     * @summary Health check
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SystemApiInterface
     */
    getHealthRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<HealthResponse>>;

    /**
     * Check if the daemon is running and healthy
     * Health check
     */
    getHealth(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<HealthResponse>;

}

/**
 *
 */
export class SystemApi extends runtime.BaseAPI implements SystemApiInterface {

    /**
     * Check if the daemon is running and healthy
     * Health check
     */
    async getHealthRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<HealthResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};


        let urlPath = `/health`;

        const response = await this.request({
            path: urlPath,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => HealthResponseFromJSON(jsonValue));
    }

    /**
     * Check if the daemon is running and healthy
     * Health check
     */
    async getHealth(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<HealthResponse> {
        const response = await this.getHealthRaw(initOverrides);
        return await response.value();
    }

}
