export interface IYChart {
    annotations: IAnnotation[];
    services: IServices;
    operations: IOperations;
}

export interface IOperations {
    [key: string]: IOperation;
}

export interface IOperation {
    [key: string]: IEdge;
}

export interface IEdge {
    name: string;
    from: string;
    to: string;
    count: number;
}

export interface IAnnotation {
    services: string[];
    operations: string[] | null;
    message: string;
    annotationType: string;
    yChartLevel: string;
    recommendation?: IRecommendation;
}

export interface IRecommendation {
    message: string;
}

export interface IServices {
    [key: string]: IService;
}

export interface IService {
    name: string;
    dependents: string[] | null;
    dependencies: string[] | null;
    cpu: IUtils;
    memory: IUtils;
    network: IUtils;
}

export interface IUtils {
    quantile: number;
    mean: number;
    stdev: number;
}
